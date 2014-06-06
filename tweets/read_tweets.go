package main

import(
  "fmt"
  "github.com/ChimeraCoder/anaconda"
  "net/url"
  "regexp"
)

const(
  ACCESS_TOKEN    = ""
  ACCESS_SECRET   = ""
  CONSUMER_KEY    = ""
  CONSUMER_SECRET = ""
)

func main(){
  anaconda.SetConsumerKey( CONSUMER_KEY )
  anaconda.SetConsumerSecret( CONSUMER_SECRET )

  api := anaconda.NewTwitterApi( ACCESS_TOKEN , ACCESS_SECRET )
  
  config := url.Values{}
  config.Set( "count", "2" )

  search , _  := api.GetSearch("mexico", config )

  for _ , tweet := range search{

    fmt.Println( tweet.Id )
    fmt.Println( tweet.User.ScreenName )
    fmt.Println( tweet.User.Name )
    fmt.Println( tweet.User.ProfileImageUrlHttps )
    fmt.Println( tweet.User.FollowersCount )
    fmt.Println( tweet.User.Lang )
    fmt.Println( tweet.User.Location )
    fmt.Println( tweet.InReplyToScreenName )
    fmt.Println( tweet.Text )

    reg := regexp.MustCompile("<a[^>]*>(.*?)</a>")
    source := reg.FindAllStringSubmatch( tweet.Source , 1 )
    fmt.Println( source[0][1] )

    for _,media := range tweet.Entities.Media{
      fmt.Println( media.Media_url_https )
    }

    for _,hashtag := range tweet.Entities.Hashtags{
      fmt.Println( hashtag.Text )
    }

    fmt.Println( tweet.Coordinates )
    fmt.Println( "\n")
  }
  
}