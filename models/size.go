package minecraft

// DigiUnit is used for performing
// storage calculations using constants.
type DigiUnit uint

const (
	Bit      DigiUnit = 1 << iota // 1
	_                             // 2
	_                             // 4
	Byte                          // 8
	Short                         // 16
	Int                           // 32
	Long                          // 64
	_                             // 128
	_                             // 256
	_                             // 512
	Kibibyte                      // 1024
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	Mibibyte                      // 1048576
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	_                             //
	Gibibyte                      // 1073741824
)
