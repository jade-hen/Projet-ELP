module Words exposing (chooseRandomIndex, loadWords, parseWords)

import Http
import Random
import String
import Types exposing (Msg(..))


loadWords : Cmd Types.Msg
loadWords =
    Http.get
        { url = "../static/words.txt", expect = Http.expectString GotWords}


-- Fichier: "word1 word2 word3 ..." (espaces / retours ligne ok)
parseWords : String -> List String
parseWords txt =
    let
        ws =
            String.words txt

        trimmed =
            List.map String.trim ws

        nonEmpty =
            List.filter (\w -> w /= "") trimmed
    in
    nonEmpty


chooseRandomIndex : List String -> Cmd Types.Msg
chooseRandomIndex ws =
    let
        n =
            List.length ws
    in
    if n <= 0 then
        Cmd.none
    else
        Random.generate PickedIndex (Random.int 0 (n - 1))
