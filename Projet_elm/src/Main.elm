module Main exposing (main)

import Browser
import Dictionary
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Types exposing (Meaning, Model, Msg(..), Status(..))
import Words


main : Program () Model Msg
main =
    Browser.element
        { init = \_ -> init
        , update = update
        , view = view
        , subscriptions = \_ -> Sub.none
        }


init : ( Model, Cmd Msg )
init =
    ( Types.emptyModel, Words.loadWords )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotWords result ->
            case result of
                Ok txt ->
                    let
                        ws =
                            Words.parseWords txt
                    in
                    ( { model | words = ws, error = "", status = Loading }
                    , Words.chooseRandomIndex ws
                    )

                Err e ->
                    ( { model | status = Error, error = Types.httpErrorToString e }
                    , Cmd.none
                    )

        PickedIndex i ->
            case Types.getAt i model.words of
                Nothing ->
                    ( { model | status = Error, error = "Liste vide / index invalide." }
                    , Cmd.none
                    )

                Just w ->
                    ( { model
                        | target = Just w
                        , meanings = []
                        , guess = ""
                        , error = ""
                        , status = Loading
                        , showSolution = False
                      }
                    , Dictionary.fetchMeanings w
                    )

        GotDefs result ->
            case result of
                Ok meanings ->
                    ( { model | meanings = meanings, status = Ready, error = "" }
                    , Cmd.none
                    )

                Err e ->
                    ( { model | status = Error, error = Types.httpErrorToString e }
                    , Cmd.none
                    )

        GuessChanged s ->
            let
                newModel =
                    { model | guess = s }

                win =
                    case model.target of
                        Nothing ->
                            False

                        Just w ->
                            Types.normalize s == Types.normalize w
            in
            if win then
                ( { newModel | status = Won }, Cmd.none )
            else
                ( newModel, Cmd.none )

        ToggleSolution ->
            ( { model | showSolution = not model.showSolution }, Cmd.none )

        NewGame ->
            ( { model
                | meanings = []
                , guess = ""
                , error = ""
                , status = Loading
                , showSolution = False
              }
            , Words.chooseRandomIndex model.words
            )


view : Model -> Html Msg
view model =
    div
        [ style "max-width" "820px"
        , style "margin" "24px auto"
        , style "font-family" "system-ui"
        , style "line-height" "1.4"
        ]
        [ h1 [] [ text "GuessIt" ]
        , p [] [ text "Lis les définitions (par catégorie) et devine le mot." ]
        , viewError model
        , viewMeanings model
        , div [ style "margin-top" "14px" ]
            [ input
                [ type_ "text"
                , value model.guess
                , onInput GuessChanged
                , placeholder "Ta réponse..."
                , style "width" "100%"
                , style "padding" "10px"
                ]
                []
            ]
        , if model.status == Won then
            p [ style "margin-top" "10px", style "font-weight" "700" ]
                [ text "✅ Correct !" ]

          else
            text ""
        , div [ style "margin-top" "12px", style "display" "flex", style "gap" "10px" ]
            [ button
                [ onClick NewGame
                , disabled (List.isEmpty model.words)
                , style "padding" "10px 14px"
                ]
                [ text "New game" ]
            , button
                [ onClick ToggleSolution
                , disabled (model.target == Nothing)
                , style "padding" "10px 14px"
                ]
                [ text (if model.showSolution then "Hide solution" else "Show solution") ]
            ]
        , viewSolution model
        ]


viewMeanings : Model -> Html msg
viewMeanings model =
    if List.isEmpty model.meanings then
        p [] [ text "Définitions : (en attente…)" ]
    else
        div []
            (List.map viewMeaningCard model.meanings)


viewMeaningCard : Meaning -> Html msg
viewMeaningCard meaning =
    div
        [ style "border" "1px solid #ddd"
        , style "border-radius" "10px"
        , style "padding" "12px 14px"
        , style "margin-top" "10px"
        ]
        [ div [ style "font-weight" "700", style "margin-bottom" "8px" ]
            [ text (meaning.partOfSpeech ++ "") ]
        , ol [ style "margin" "0 0 0 18px" ]
            (List.map (\d -> li [ style "margin-bottom" "6px" ] [ text d ]) meaning.definitions)
        ]


viewSolution : Model -> Html msg
viewSolution model =
    if model.showSolution then
        case model.target of
            Nothing ->
                text ""

            Just w ->
                p [ style "margin-top" "10px", style "font-weight" "700" ]
                    [ text ("Solution : " ++ w) ]
    else
        text ""


viewError : Model -> Html msg
viewError model =
    if model.error == "" then
        text ""
    else
        p [ style "color" "crimson" ] [ text ("Erreur : " ++ model.error) ]
