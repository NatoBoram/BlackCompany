package nord

// Nord consists of four named color palettes providing different syntactic
// meanings and color effects for dark & bright ambiance designs.
type Colour string

// Polar Night is made up of four darker colors that are commonly used forbase
// elements like backgrounds or text color in bright ambiance designs.
const (
	Nord00 Colour = "#2e3440"
	Nord01 Colour = "#3b4252"
	Nord02 Colour = "#434c5e"
	Nord03 Colour = "#4c566a"
)

// Snow Storm is made up of three bright colors that are commonly used for text
// colors or base UI elements in bright ambiance designs.
const (
	Nord04 Colour = "#d8dee9"
	Nord05 Colour = "#e5e9f0"
	Nord06 Colour = "#eceff4"
)

// Frost can be described as the heart palette of Nord, a group of four bluish
// colors that are commonly used for primary UI component and text highlighting
// and essential code syntax elements.
const (
	Nord07 Colour = "#8fbcbb"
	Nord08 Colour = "#88c0d0"
	Nord09 Colour = "#81a1c1"
	Nord10 Colour = "#5e81ac"
)

// Aurora consists of five colorful components reminiscent of the "Aurora
// borealis", sometimes referred to as polar lights or northern lights.
const (
	Nord11 Colour = "#bf616a"
	Nord12 Colour = "#d08770"
	Nord13 Colour = "#ebcb8b"
	Nord14 Colour = "#a3be8c"
	Nord15 Colour = "#b48ead"
)

func (c Colour) String() string {
	return string(c)
}
