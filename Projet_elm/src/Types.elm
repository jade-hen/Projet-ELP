module Types exposing
    ( Model, Msg(..), Status(..), Meaning, emptyModel, normalize, getAt, httpErrorToString)

import Http
import String


type Status
    = Loading| Ready| Won| Error


type alias Meaning =
    { partOfSpeech : String, definitions : List String}


type alias Model =
    { words : List String, target : Maybe String, meanings : List Meaning, guess : String, status : Status, error : String, showSolution : Bool}


emptyModel : Model
emptyModel =
    { words = [], target = Nothing, meanings = [], guess = "", status = Loading, error = "", showSolution = False}


type Msg
    = GotWords (Result Http.Error String)| PickedIndex Int| GotDefs (Result Http.Error (List Meaning))| GuessChanged String| NewGame| ToggleSolution


normalize : String -> String
normalize s =
    String.toLower (String.trim s)


getAt : Int -> List a -> Maybe a
getAt i xs =
    if i < 0 then
        Nothing
    else
        case ( i, xs ) of
            ( 0, y :: _ ) ->
                Just y

            ( n, _ :: rest ) ->
                getAt (n - 1) rest

            ( _, [] ) ->
                Nothing


httpErrorToString : Http.Error -> String
httpErrorToString err =
    case err of
        Http.BadUrl u ->
            "BadUrl: " ++ u

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "NetworkError"

        Http.BadStatus code ->
            "BadStatus: " ++ String.fromInt code

        Http.BadBody msg ->
            "BadBody: " ++ msg
