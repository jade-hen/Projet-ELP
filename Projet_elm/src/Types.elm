module Types exposing
    ( Model, Msg(..), Status(..), Meaning, emptyModel, normalize, getAt, httpErrorToString)

import Http
import String


type Status --statut du jeu en cours
    = Loading | Ready | Won | Error


type alias Meaning = -- pour garder la liste des définitions
    { partOfSpeech : String, definitions : List String}


type alias Model = -- le modèle qui garde l'affichage en cours
    { words : List String, target : Maybe String, meanings : List Meaning, guess : String, status : Status, error : String, showSolution : Bool, points : Int, solutionShown : Bool}


emptyModel : Model -- fonction pour vider le modèle (revenir à la page de départ)
emptyModel =
    { words = [], target = Nothing, meanings = [], guess = "", status = Loading, error = "", showSolution = False, points = 0, solutionShown = False }


type Msg -- msg à envoyer pour changer la vue (qd le joueur modifie un paramètre)
    = GotWords (Result Http.Error String) | PickedIndex Int | GotDefs (Result Http.Error (List Meaning)) | GuessChanged String | NewGame | ToggleSolution


normalize : String -> String -- fonction pour mettre en minuscule 
normalize s =
    String.toLower (String.trim s)


getAt : Int -> List a -> Maybe a -- récupérer le mot à l'indice i de la liste de mots de façon récursive
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


httpErrorToString : Http.Error -> String -- pour afficher les erreurs de requêtes http
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
