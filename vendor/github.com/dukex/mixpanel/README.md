
mixpanel
========

Mixpanel Go Client

## Donate to 1P8ccYhVt4ByLahuVXiCY6U185gmYA8Rqf

## Usage

``` go
import "github.com/dukex/mixpanel"
```
--

[documentation on godoc](http://godoc.org/github.com/dukex/mixpanel)


## Examples

Track

``` go
err := client.Track("13793", "Signed Up", map[string]interface{}{
	"Referred By": "Friend",
})
```
--

Identify and Update Operation

``` go
people := client.Identify("13793")

err := people.Track(map[string]interface{}{
	"Buy": "133"
})

err := people.Update("$set", map[string]interface{}{
	"Address":  "1313 Mockingbird Lane",
	"Birthday": "1948-01-01",
})
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).

## Author

Duke X ([dukex](http://github.com/dukex))
