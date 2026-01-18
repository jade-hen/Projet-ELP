module Dictionary exposing (fetchMeanings)

import Http
import Json.Decode as D
import Types exposing (Meaning, Msg(..))


fetchMeanings : String -> Cmd Types.Msg -- la fonction prend en entrée un mot et utilise l'api pour récupérer la définition (url ++ mot)
fetchMeanings word =                    -- et ensuite ça envoie le message comme quoi on a reçu (ou pas) la def ? 
    Http.get
        { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ word, expect = Http.expectJson GotDefs responseDecoder} 


-- L'API renvoie une LISTE d'entries, on prend la première (c'est quoi les autres entries ? je crois qu'il y en a pas toujours d'autres)
responseDecoder : D.Decoder (List Meaning)
responseDecoder =
    D.andThen firstOrFail (D.list entryDecoder) -- json -> list des meanings


firstOrFail : List a -> D.Decoder a
firstOrFail xs =
    case xs of
        x :: _ ->
            D.succeed x

        [] -> --si c'est une liste vide
            D.fail "Empty response"


entryDecoder : D.Decoder (List Meaning)
entryDecoder =
    D.field "meanings" (D.list meaningDecoder) --récupère le champ meanings du json


meaningDecoder : D.Decoder Meaning -- récupère les champs partOfSpeech et definitions du json
meaningDecoder =
    D.map2 Meaning
        (D.field "partOfSpeech" D.string)
        (D.field "definitions" (D.list (D.field "definition" D.string)))
