module Dictionary exposing (fetchMeanings)

import Http
import Json.Decode as D
import Types exposing (Meaning, Msg(..))


fetchMeanings : String -> Cmd Types.Msg
fetchMeanings word =
    Http.get
        { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ word
        , expect = Http.expectJson GotDefs responseDecoder
        }


-- L'API renvoie une LISTE d'entries, on prend la première
responseDecoder : D.Decoder (List Meaning)
responseDecoder =
    D.list entryDecoder
        |> D.andThen firstOrFail


firstOrFail : List a -> D.Decoder a
firstOrFail xs =
    case xs of
        x :: _ ->
            D.succeed x

        [] ->
            D.fail "Empty response"


entryDecoder : D.Decoder (List Meaning)
entryDecoder =
    D.field "meanings" (D.list meaningDecoder)


meaningDecoder : D.Decoder Meaning
meaningDecoder =
    D.map2 Meaning
        (D.field "partOfSpeech" D.string)
        (D.field "definitions" (D.list (D.field "definition" D.string)))
