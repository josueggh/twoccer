package readtweets

import(
  "github.com/ChimeraCoder/anaconda"
  "labix.org/v2/mgo"
  "net/url"
  "regexp"
)

const(
  DATABASE  = "twoccer"
)


type Readtweets struct{
  accesstoken     string
  accessecret     string
  consumerkey     string
  consumersecret  string
}

type Tweet struct{
  Ti    int64   `json:"tid"`
  Sn    string  `json:"screename"`
  Pp    string  `json:"profilepicture"`
  Fc    int     `json:"followers_count"`
  Lg    string  `json:"lang"`
  Lt    string  `json:"location"`
  Rp    string  `json:"replyto"`
  Sc    string  `json:"source"`
  Md    []string  `bson:"media,omitempty"` 
  Ht    []string  `bson:"hashtag,omitempty"` 
  Tx    string  `json:"text"`
}

func New(access_token, access_secret , consumer_key, consumer_secret string) *Readtweets {
  r := &Readtweets{
    accesstoken     : access_token,
    accessecret     : access_secret,
    consumerkey     : consumer_key,
    consumersecret  : consumer_secret,
  }
  return r
}


func (self *Readtweets) Save( collection , words , limit string ) bool{

  session, err := mgo.Dial( "localhost" )
  if err != nil {
    panic(err)
  }

  defer session.Close()
  session.SetMode(mgo.Monotonic, true )

  tweets_collection :=session.DB( DATABASE ).C( collection )

  anaconda.SetConsumerKey( self.consumerkey )
  anaconda.SetConsumerSecret( self.consumersecret )

  api := anaconda.NewTwitterApi( self.accesstoken , self.accessecret )
  
  config := url.Values{}
  config.Set( "count", limit )

  search , _  := api.GetSearch( words , config )

  for _ , tweet := range search{

    reg          := regexp.MustCompile( "<a[^>]*>(.*?)</a>" )
    source       := reg.FindAllStringSubmatch( tweet.Source , 1 )
    real_source  := source[0][1]
  
    media_list   := [] string{}
    hashtag_list := [] string{}

    for _, media := range tweet.Entities.Media{
      media_list = append( media_list , media.Media_url_https )
    }

    for _,hashtag := range tweet.Entities.Hashtags{
      hashtag_list = append( hashtag_list , hashtag.Text )
    }

    t:= &Tweet{ 
      tweet.Id,
      tweet.User.ScreenName,
      tweet.User.ProfileImageUrlHttps,
      tweet.User.FollowersCount,
      tweet.User.Lang,
      tweet.User.Location,
      tweet.InReplyToScreenName,
      real_source,
      media_list,
      hashtag_list,
      tweet.Text,
    }

    tweets_collection.Insert( t )
    
  }
  
  return true
}