module Dictionary exposing (fetchMeanings)

import Http
import Json.Decode exposing (Decoder, map2, andThen, list, string, field, succeed, fail)
import Types exposing (Meaning, Msg(..))


fetchMeanings : String -> Cmd Types.Msg -- la fonction prend en entrée un mot et utilise l'api pour récupérer la définition (url ++ mot)
fetchMeanings word =                    -- et ensuite ça envoie le message comme quoi on a reçu (ou pas) la def  
    Http.get
        { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ word, expect = Http.expectJson GotDefs responseDecoder} 


-- L'API renvoie une LISTE d'entries, on prend la première (c'est quoi les autres entries ? je crois qu'il y en a pas toujours d'autres)
responseDecoder : Decoder (List Meaning)
responseDecoder =
    andThen firstOrFail (list entryDecoder) -- json -> list des meanings


firstOrFail : List a -> Decoder a
firstOrFail xs =
    case xs of
        x :: _ ->
            succeed x

        [] -> --si c'est une liste vide
            fail "Empty response"


entryDecoder : Decoder (List Meaning)
entryDecoder =
    field "meanings" (list meaningDecoder) --récupère le champ meanings du json


meaningDecoder : Decoder Meaning -- récupère les champs partOfSpeech et definitions du json
meaningDecoder =
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list (field "definition" string)))
