import "./Chess.css"
import { getStandrardBoard, Game } from "../game/board"
import { useState } from "react"


function Square(props) {
    return (
        <td
            className="Square"
            onClick={props.onClick}
            style={props.selected ? { border: "2px solid red" } : { border: "2px solid white" }}
        >
            <img alt="" src={props.image} />
        </td>
    )
}

const game = new Game(getStandrardBoard());

const rowNumbs = [8, 7, 6, 5, 4, 3, 2, 1]
const colNumbs = ["a", "b", "c", "d", "e", "f", "g", "h"]

export function Board() {
    const [board, setBoard] = useState(game.board);
    const [selectedPiece, setSelectedPiece] = useState(null);

    const movePiece = (evt, col, row) => {
        if (selectedPiece) {
            if (game.movePiece({ row: selectedPiece.row, col: selectedPiece.col }, { row: row, col: col })) {
                setBoard(game.board)
                setSelectedPiece(null)
            } else {
                setSelectedPiece(null)
            }
        } else if (board[row][col]) {
            setSelectedPiece({ col: col, row: row })
        }
    }

    return (
        <table className="Board">
            <tbody>
                {board.map((row, rowIndex) => {
                    return (
                        <tr>
                            <th>{rowNumbs[rowIndex]}</th>
                            {row.map((square, colIndex) => {
                                if (square) {
                                    return (
                                        <Square
                                            onClick={(evt, col, row) => { movePiece(evt, colIndex, rowIndex) }}
                                            image={square.image}
                                            selected={(selectedPiece?.col === colIndex && selectedPiece?.row === rowIndex) ? true : null}
                                        />
                                    )
                                } else {
                                    return (
                                        <Square
                                            onClick={(evt, col, row) => { movePiece(evt, colIndex, rowIndex) }}
                                        />
                                    )
                                }
                            })}
                        </tr>
                    )
                })}
                <tr>
                    <th></th>
                    {board.map((row, rowIndex) => {
                        return <th>{colNumbs[rowIndex]}</th>
                    })}
                </tr>
            </tbody>
        </table>
    )
}