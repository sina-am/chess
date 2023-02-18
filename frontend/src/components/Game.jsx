
import { Box, Button, FormControl, Select, InputLabel, MenuItem } from "@mui/material";
import { getTokenFromLocalStorage } from "../lib/auth";
import { Board } from "./Chess";

import React, { useState, useCallback, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

// export const GameOptions = () => {
//     const socketUrl = 'ws://localhost:8080/ws?token=' + getTokenFromLocalStorage();
//     const [messageHistory, setMessageHistory] = useState([]);
//     const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

//     useEffect(() => {
//         if (lastMessage !== null) {
//             setMessageHistory((prev) => prev.concat(lastMessage));
//         }
//     }, [lastMessage, setMessageHistory]);

//     const handleClickSendMessage = useCallback(() => sendMessage(), []);

//     const connectionStatus = {
//         [ReadyState.CONNECTING]: 'Connecting',
//         [ReadyState.OPEN]: 'Open',
//         [ReadyState.CLOSING]: 'Closing',
//         [ReadyState.CLOSED]: 'Closed',
//         [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
//     }[readyState];

//     return (
//         <div>
//             <button
//                 onClick={handleClickSendMessage}
//                 disabled={readyState !== ReadyState.OPEN}
//             >
//                 Click Me to send 'Hello'
//             </button>
//             <span>The WebSocket is currently {connectionStatus}</span>
//             {lastMessage ? <span>Last message: {lastMessage.data}</span> : null}
//             <ul>
//                 {messageHistory.map((message, idx) => (
//                     <span key={idx}>{message ? message.data : null}</span>
//                 ))}
//             </ul>
//         </div>
//     );
// };


export function Game() {
    return (
        <Board />
    )
}