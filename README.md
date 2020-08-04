# go-validate

This package aims to provide a few validations for your request objects using Tag approach.

## Usage
```
...

type MyRequest struct {
  ID        string  `validate:"required"`                 // Set ID as a required field.
  Email     string  `validate:"required,email,max=255"`   // Set Email as required and check if it is in the Email format
  Password  string  `validate:"required,password"`        // Set Password as required and check if has a strong format: at least one Uppercase, Special Character,    Number and if has a length between 8-20.
  Age       int     `validate:"min=18,max=64"`            // Check if the number in the range.
}

func main() {
  request := MyRequest{"123", "test@mail.com", "StrongP@ssw0rd", 20}
  
  fieldErrors := govalidate.BodyRequest(d)
	if fieldErrors != nil {
		return &domainError.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid body fields",
			InError: fieldErrors,
		}
	}
}
```

For each invalid field the validator will return a string telling what field is wrong and why.
