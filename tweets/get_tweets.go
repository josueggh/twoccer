package main

import(
"fmt"
"github.com/blacknash/twoccer/tweets/v1"
"labix.org/v2/mgo"
"labix.org/v2/mgo/bson"
 "time"
 "os"
)

const(
  ACCESS_TOKEN    = "YOUR ACCESS TOKEN"
  ACCESS_SECRET   = "YOUR ACCESS SECRET"
  CONSUMER_KEY    = "YOUR CONSUMER KEY"
  CONSUMER_SECRET = "YOUR CONSUMER SECRET"
  DATABASE        = "twoccer"
  MINUTES         = 180
  UPDATE_TIME     = 120
  TWEETS          = "100"
)

type empty struct{}

type Match struct{
  Id bson.ObjectId `bson:"_id,omitempty"`
  Date  time.Time
  TeamA string
  TeamB string
}

type Team struct{
  Words []string
}

func main(){

  fmt.Println("searching matches")

  session, err := mgo.Dial( "localhost" )
  if err != nil {
    panic(err)
    os.Exit(1)
  }

  defer session.Close()
  session.SetMode(mgo.Monotonic, true )

  matches :=  session.DB( DATABASE ).C( "matchestest" )
  teams   :=  session.DB( DATABASE ).C( "teams" )

  result:= Match{}
  teamA := Team{}
  teamB := Team{}
  
  errM := matches.Find( 
            bson.M{ 
                "date"    : bson.M{ "$lt" : time.Now() } , 
                "status"  : bson.M{ "$nin": bson.M{ "finished": "finished" , "started" : "started"} } ,
          }).One(&result)

  if errM !=nil {
    os.Exit(1)
  }

  fmt.Println("Id", result.Id)

  match := result.TeamA+result.TeamB
  
  errTA := teams.Find(
    bson.M{
      "_id" : result.TeamA,
  }).One(&teamA)

  errTB := teams.Find(
    bson.M{
      "_id" : result.TeamB,
  }).One(&teamB)

  if errTA !=nil && errTB!=nil{
    os.Exit(1)
  }

  words := []string{"gol","gool","goool","gooool","goooool","gooooool","gooooool","gooooool","gooooool","goooooool","gooooooool","goooooooool","winner","loser","#GOL"}

  for _ , word := range teamA.Words{
    words = append( words , word )
  }

  for _ , word := range teamB.Words{
    words = append( words , word )
  }

  reader := readtweets.New(ACCESS_TOKEN, ACCESS_SECRET, CONSUMER_KEY , CONSUMER_SECRET)

  ticker := time.NewTicker(time.Second * UPDATE_TIME)

  go func() {
      for t := range ticker.C {

        for _, search := range words{
          fmt.Println( search )
          reader.Save(  DATABASE , match  , search , TWEETS )
        }
        fmt.Println("Tick at", t)
      }
  }()

  time.Sleep(time.Minute * MINUTES)
  
  ticker.Stop()
  fmt.Println("Match finished")

}