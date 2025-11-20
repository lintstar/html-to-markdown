package escape

import "github.com/lintstar/html-to-markdown/v2/marker"

var placeholderRune rune = marker.MarkerEscaping

// IMPORTANT: Only internally we assume it is only byte
var placeholderByte byte = marker.BytesMarkerEscaping[0]
