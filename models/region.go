package minecraft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/KellyLSB/nbt"
	"github.com/KellyLSB/nbt/alt/minecraft/mmap"
	"github.com/KellyLSB/utils/errors"
	"github.com/davecgh/go-spew/spew"
)

type RegionPos int32

func (rp *RegionPos) ReadFrom(r io.Reader) (int64, error) {
	if err := binary.Read(r, binary.BigEndian, rp); err != nil {
		return 0, err
	}

	return 4, nil
}

func (rp *RegionPos) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, binary.BigEndian, rp); err != nil {
		return 0, err
	}

	return 4, nil
}

func (rp *RegionPos) SetLocation(location int32) {
	*rp = RegionPos((location << 8) | int32(int8(*rp)))
}

func (rp *RegionPos) Location() int32 {
	return int32(*rp) >> 8
}

func (rp *RegionPos) SetSectors(sectors int8) {
	*rp = RegionPos(((int32(*rp) >> 8) << 8) | int32(sectors))
}

func (rp *RegionPos) Sectors() int8 {
	return int8(*rp)
}

func (rp *RegionPos) NotExists() bool {
	return int32(*rp) < 1
}

// RegionMod is the last modification time of a chunk. Unit: unknown, seconds?
//
// NOTE: Does something use this? MCEdit maybe?
type RegionMod int32

func (rm *RegionMod) ReadFrom(r io.Reader) (int64, error) {
	if err := binary.Read(r, binary.BigEndian, rm); err != nil {
		return 0, err
	}

	return 4, nil
}

func (rm *RegionMod) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, binary.BigEndian, rm); err != nil {
		return 0, err
	}

	return 4, nil
}

type RegionData struct {
	Length      int32
	Compression mmap.Compression
	Chunk       Root
}

func (rd *RegionData) ReadFrom(r io.Reader) (n int64, err error) {
	// Read Length
	if err := binary.Read(r, binary.BigEndian, &rd.Length); err != nil {
		return n, err
	} else {
		n += 4
	}

	// Read Compression
	if err := binary.Read(r, binary.BigEndian, &rd.Compression); err != nil {
		return n, err
	} else {
		n++
	}

	// Read Chunk
	var buf bytes.Buffer
	buf.ReadFrom(mmap.ArchiveReader(rd.Compression, r))
	if err := nbt.NewDecoder(&buf).Decode(&rd.Chunk); err != nil {
		return n, err
	} else {
		n += int64(rd.Length)
	}

	return n, nil
}

func (rd *RegionData) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	// Write Chunk to Buffer
	// if err := nbt.NewEncoder(
	// 	mmap.ArchiveWriter(rd.Compression, buf),
	// ).Encode(&rd.Chunk); err != nil {
	// 	return n, err
	// } else {
	// 	rd.Length = buf.Len() + 1
	// }

	// Write Length
	if err := binary.Write(w, binary.BigEndian, &rd.Length); err != nil {
		return n, err
	} else {
		n += 4
	}

	// Write Compression
	if err := binary.Write(w, binary.BigEndian, &rd.Compression); err != nil {
		return n, err
	} else {
		n++
	}

	// Write Buffer
	n1, err := buf.WriteTo(w)
	return n + n1, err
}

// Region (.mca files) store 32x32 chunks.
type Region struct {
	Pos  [1024]RegionPos // Chunk position in 4k increments from start.
	Mod  [1024]RegionMod // Last modification time of a chunk.
	Data [1024]*RegionData
}

func (re *Region) ReadFrom(r io.Reader) (n int64, err error) {
	var (
		buf      = mmap.ReadMMap(r)
		errs     errors.Errors
		totalLen int64
	)

	errs.AllOrNothing()

	spew.Printf(
		"==> Bytes Read %04d (%04x); Beginning chunk position decoding.\n",
		totalLen, totalLen,
	)

	// Read chunk positions.
	for i := 0; i < len(re.Pos); i++ {
		errs.Step(func(err *error) {
			n, *err = re.Pos[i].ReadFrom(buf)
			totalLen += n
		})
	}

	spew.Printf(
		"==> Bytes Read %04d (%04x); Beginning chunk timestamp decoding.\n",
		totalLen, totalLen,
	)

	// Read chunk timestamps.
	for i := 0; i < len(re.Mod); i++ {
		errs.Step(func(err *error) {
			n, *err = re.Mod[i].ReadFrom(buf)
			totalLen += n
		})
	}

	spew.Printf(
		"==> Bytes Read %04d (%04x); Beginning chunk data decoding.\n",
		totalLen, totalLen,
	)

	// Read chunks from the buffer.
	for i := 0; i < len(re.Pos); i++ {
		if re.Pos[i].NotExists() {
			continue
		}

		offsetK := re.Pos[i].Location()
		lengthK := re.Pos[i].Sectors()
		offsetB := int64(offsetK) * 4096
		lengthB := int64(lengthK) * int64(Kibibyte)

		spew.Printf(
			"==> Chunk %d of 1024 is located at offset %dKib (%dB) "+
				"with a length of %dKib (%dB).\n",
			i, offsetK, offsetB, lengthK, lengthB,
		)

		if re.Data[i] == nil {
			re.Data[i] = new(RegionData)
		}

		errs.Step(func(err *error) {
			n, *err = re.Data[i].ReadFrom(buf.Range(offsetB, lengthB))
			totalLen += n
		})
	}

	return totalLen, fmt.Errorf(errs.Error())
}

func (re *Region) WriteTo(w io.Writer) (n int64, err error) {
	var (
		buf  = mmap.MakeMMap(0).AppendModeOn()
		errs errors.Errors
	)

	errs.AllOrNothing()

	spew.Printf(
		"==> Bytes Written %04d (%04x); Beginning chunk position encoding.\n",
		n, n,
	)

	// Write chunk positions.
	for i := 0; i < len(re.Pos); i++ {
		errs.Step(func(err *error) {
			_, *err = re.Pos[i].WriteTo(buf)
		})
	}

	spew.Printf(
		"==> Bytes Written %04d (%04x); Beginning chunk timestamp encoding.\n",
		n, n,
	)

	// Write chunk timestamps.
	for i := 0; i < len(re.Mod); i++ {
		errs.Step(func(err *error) {
			_, *err = re.Mod[i].WriteTo(buf)
		})
	}

	spew.Printf(
		"==> Bytes Read %04d (%04x); Beginning chunk data decoding.\n",
		n, n,
	)

	// Write chunks from the buffer.
	for i := 0; i < len(re.Pos); i++ {
		if re.Pos[i].NotExists() {
			continue
		}

		offsetK := re.Pos[i].Location()
		lengthK := re.Pos[i].Sectors()
		offsetB := int64(offsetK) * 4096
		lengthB := int64(lengthK) * int64(Kibibyte)

		spew.Printf(
			"==> Chunk %d of 1024 is located at offset %dKib (%dB) "+
				"with a length of %dKib (%dB).\n",
			i, offsetK, offsetB, lengthK, lengthB,
		)

		if re.Data[i] == nil {
			re.Data[i] = new(RegionData)
		}

		errs.Step(func(err *error) {
			_, *err = re.Data[i].WriteTo(buf.Range(offsetB, lengthB))
		})
	}

	errs.Step(func(err *error) {
		n, *err = buf.WriteTo(w)
	})

	return n, fmt.Errorf(errs.Error())
}

func (re *Region) ChunkPos(x, z int32) int { return int(z<<5 + x) }
