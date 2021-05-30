# GABAHTMLParser

A Simple HTML Parser.

## Installation
```shell
go get github.com/HarrisonKawagoe3960X/GABAHTMLParser
```
## Usage
### Import the GABAHTMLParser
Just add ```"github.com/HarrisonKawagoe3960X/GABAHTMLParser"``` into ```import``` like this:
```go
package main

import(
	"fmt"
	"github.com/HarrisonKawagoe3960X/GABAHTMLParser" //Add this
)

func main() {
	htmlobject := GABAHTMLParser.GetHTMLfromURL("someurl",false)
	results := htmlobject.Find("tag = 'a'")
	for _ , result := range results{
		fmt.Println(result.InnerHTML)
	}
	
}
```

### Parse HTML from URL or path
```go
htmlobject := GABAHTMLParser.GetHTMLfromURL("someurl",false)
```
if you parse the source of site that use Shift-JIS encoding, change ```false``` to ```true```
```go
htmlobject := GABAHTMLParser.GetHTMLfromURL("someurl",true)
```

### Parse HTML from String Array
```go
htmlobject := GABAHTMLParser.ParseHTML(strarray)
```


### Element Object
After parsing the HTML, you can extract the data by calling ```Element```. 
- ```InnerHTML```: HTML code under the current HTML Element.
- ```Tag```: Tag name of current HTML Element.
-  ```Child```: Child Objects of current HTML Element.
-  ```Parent```: Parent Object of current HTML Element.
-  ```Attr```: Properties of current HTML Element.
### Search for HTML Object 
```go
results := htmlobject.Find("tag = 'a'")
```
you can combine conditions by using ```&&```
```go
results := htmlobject.Find("tag = 'a' && class = 'hanshin'")
```

