# removebg
A Golang API wrapper for removing background using [remove.bg](https://www.remove.bg/)'s [API]()  Topics

# License
This code is licensed under the GPL3 License. See [here](https://github.com/g-lib/removebg/blob/main/LICENSE) for more details.

# Installation
`go get github.com/g-lib/removebg@latest`

# Usage

## `RemoveFromFile`

Removes the background given an image file.



### Code Example:
```go
import "github.com/g-lib/removebg"
rmbg = NewRemoveBg("YOUR-API-KEY")
rmbg.RemoveFromFile("test.jpg",nil)
```


## `RemoveBackgroundFromImgURL`

Removes the background given an image URL.



### Code Example:
```go
import "github.com/g-lib/removebg"
rmbg = NewRemoveBg("YOUR-API-KEY")
rmbg.RemoveFromURL("http://www.example.com/some_image.jpg",nil)
```


## `RemoveBackgroundFromBase64Img`

Removes the background given a base64 image string.



### Code Example:
```go
import "github.com/g-lib/removebg"
rmbg = NewRemoveBg("YOUR-API-KEY")
rmbg.RemoveFrombBase64("BASE64-CONTENT",nil)
```

