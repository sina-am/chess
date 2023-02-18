import { Dashboard } from "../components/layouts/Dashboard";
import { Game } from "../components/Game"

import { Box, Button, FormControl, Select, InputLabel, MenuItem } from "@mui/material";
import { getTokenFromLocalStorage } from "../lib/auth";
import { Board } from "../components/Chess";

import React, { useState, useCallback, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';


export function GamePage() {
    const [step, setStep] = useState(0);
    const [option, setOption] = useState("");
    const [messageHistory, setMessageHistory] = useState([]);
    const [duration, setDuration] = useState(10);
    const { sendMessage, lastMessage, readyState } = useWebSocket('ws://localhost:8080/ws?token=' + getTokenFromLocalStorage());

    useEffect(() => {
        if (lastMessage !== null) {
            setMessageHistory((prev) => prev.concat(lastMessage));
        }
    }, [lastMessage, setMessageHistory]);

    const handleOnlineGameRequest = useCallback(() => {
        setStep(1)
        setOption("ONLINE")
        sendMessage(JSON.stringfy({
            type: "FRINDLY",
            duration: parseInt(duration, 10),
        }))
    }, []);

    if (step === 0) {
        return (
            <Dashboard>
                <Box
                    sx={{
                        marginLeft: "10%",
                        marginTop: "10%",
                        maxWidth: "500px",
                        '& .MuiTextField-root': {
                            display: "flex",
                            marginTop: "20px",
                        },

                    }}
                >
                    <FormControl fullWidth>
                        <InputLabel id="select game duration">Game Duration</InputLabel>
                        <Select
                            labelId="demo-simple-select-label"
                            id="demo-simple-select"
                            value={duration}
                            label="Duration"
                            onChange={(evt) => { setDuration(evt.target.value) }}
                        >
                            <MenuItem value={10}>Ten</MenuItem>
                            <MenuItem value={20}>Twenty</MenuItem>
                            <MenuItem value={30}>Thirty</MenuItem>
                        </Select>
                    </FormControl>

                    <Button
                        fullWidth
                        style={{ marginTop: "20px" }}
                        variant="contained"
                        onClick={handleOnlineGameRequest}
                    >Online game</Button>

                    <Button
                        fullWidth
                        style={{ marginTop: "20px" }}
                        variant="contained"
                        onClick={() => {
                            setStep(1)
                            setOption("FRIENDS")
                        }}
                    > Play a friend </Button>

                    <Button
                        fullWidth
                        style={{ marginTop: "20px" }}
                        variant="contained"
                        onClick={() => {
                            setStep(1)
                            setOption("SINGLE")
                        }}
                    > Vs computer </Button>
                </Box>
            </Dashboard>
        )
    }

    else if (step === 1) {
        switch (option) {
            case "SINGLE":
                return (
                    <Dashboard>
                        <Game />
                    </Dashboard>
                )

            case "ONLINE":
                
                return (
                    <Dashboard>
                        {lastMessage ? <span>Last message: {lastMessage.data}</span> : null}
                        <ul>
                            {messageHistory.map((message, idx) => (
                                <span key={idx}>{message ? message.data : null}</span>
                            ))}
                        </ul>
                        <Game />
                    </Dashboard>
                )
            default:
                return (
                    <Dashboard>
                        <Game />
                    </Dashboard>
                )
        }
    }
}