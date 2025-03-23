package main

type S3Event struct {
	Event   string   `json:"Event"`
	Records []Record `json:"Records"`
}

type Record struct {
	S3 S3 `json:"s3"`
}

type S3 struct {
	Bucket Bucket `json:"bucket"`
	Object Object `json:"object"`
}

type Bucket struct {
	Name string `json:"name"`
}

type Object struct {
	Key string `json:"key"`
}
