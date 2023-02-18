import { Box, TextField, Button, FormControlLabel, Radio, FormControl, FormLabel, RadioGroup, Autocomplete } from "@mui/material";
import { useState, useContext } from "react";
import { UserContext } from "../App";
import { useUser } from '../lib/auth'


export function Profile() {
    const {user, _} = useContext(UserContext);
    const countries = [
        "Germany", "America"
    ]
    const [password, setPassword] = useState("");

    if (!user) {
        return;
    }

    return (
        <Box
            sx={{
                marginLeft: "10%",
                marginTop: "10%",
                '& .MuiTextField-root': {
                    maxWidth: "500px",
                    display: "flex",
                    marginTop: "20px",
                },

            }}
        >
            <TextField id="email-input" label="Email" variant="filled" value={user.email} disabled/>
            <TextField
                variant="filled"
                label="Password"
                value={password}
                type="password"
                onChange={(e) => { setPassword(e.target.value); }}
            />
            <Autocomplete
                disablePortal
                id="combo-box-demo"
                options={countries}
                renderInput={(params) => <TextField {...params} label="Country" />}
            />
            <FormControl style={{display: "flex", marginTop: "20px"}}>
                <FormLabel id="demo-row-radio-buttons-group-label">Gender</FormLabel>
                <RadioGroup
                    row
                    aria-labelledby="demo-row-radio-buttons-group-label"
                    name="row-radio-buttons-group"
                >
                    <FormControlLabel value="female" control={<Radio />} label="Female" />
                    <FormControlLabel value="male" control={<Radio />} label="Male" />
                    <FormControlLabel value="other" control={<Radio />} label="Other" />
                </RadioGroup>
            </FormControl>
            

            <Button
                style={{ marginTop: "20px" }}
                variant="contained"
            >
                Save
            </Button>
        </Box>
    );
}