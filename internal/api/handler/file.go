package handler

// File struct represents a file returned as response
type File struct {
	Name        string `json:Name`
	Extension   string `json:Extension`
	SizeInBytes int    `json:SizeInBytes`
	FileLink    string `json:FileLink`
}
