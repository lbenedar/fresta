package db

type Wall struct {
	ID string

	C         []int
	Light     int
	Move      int
	Sight     int
	Sound     int
	Dir       int
	Door      int
	Ds        int
	Threshold Threshold
	Animation any
}

type Threshold struct {
	ID uint

	Light       int
	Sight       int
	Sound       int
	Attenuation bool
}
