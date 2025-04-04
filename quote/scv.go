package quote

// SCV Trained
const SCVTrained = "SCV ready."

// When attacked
const (
	SCVAttacked1 = "I'm too young to die!"
	SCVAttacked2 = "Help!"
	SCVAttacked3 = "Not what I had in mind!"
)

// Selected
const (
	SCVSelected1 = "Huh?"
	SCVSelected2 = "What's goin' on?"
	SCVSelected3 = "Bad news?"
	SCVSelected4 = "Aah! Ya scared me!"
	SCVSelected5 = "Go ahead."
	SCVSelected6 = "Big job, huh?"
	SCVSelected7 = "In the rear with the gear."
)

// Move order
const (
	SCVMove01 = "You're the boss."
	SCVMove02 = "Yep."
	SCVMove03 = "Yeah whatever."
	SCVMove04 = "It's your dime."
	SCVMove05 = "Woo hoo! Overtime!"
	SCVMove06 = "Well butter my biscuit!"
	SCVMove07 = "Yeah."
	SCVMove08 = "Yup."
	SCVMove09 = "Yo."
	SCVMove10 = "Gotcha."
	SCVMove11 = "Uh-huh..."
	SCVMove12 = "Will do."
	SCVMove13 = "We hear ya."
	SCVMove14 = "Yes sir!"
	SCVMove15 = "Sure thing."
	SCVMove16 = "Roger."
	SCVMove17 = "I'm goin'!"
	SCVMove18 = "Move it!"
	SCVMove19 = "Yeah, yeah."
)

// Attack order
const (
	SCVAttack1 = "This is crazy!"
	SCVAttack2 = "This is your plan!?"
	SCVAttack3 = "What, you run out of marines?"
	SCVAttack4 = "Oh that's just great..."
)

// SCVQuotes contains all the SCV quotes
var SCVQuotes = []string{
	SCVTrained,

	SCVAttacked1, SCVAttacked2, SCVAttacked3,

	SCVSelected1, SCVSelected2, SCVSelected3, SCVSelected4, SCVSelected5,
	SCVSelected6, SCVSelected7,

	SCVMove01, SCVMove02, SCVMove03, SCVMove04, SCVMove05, SCVMove06, SCVMove07,
	SCVMove08, SCVMove09, SCVMove10, SCVMove11, SCVMove12, SCVMove13, SCVMove14,
	SCVMove15, SCVMove16, SCVMove17, SCVMove18, SCVMove19,

	SCVAttack1, SCVAttack2, SCVAttack3, SCVAttack4,
}
