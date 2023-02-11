import * as React from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Divider from '@mui/material/Divider';
import ListItemText from '@mui/material/ListItemText';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';
import Typography from '@mui/material/Typography';
import axios from 'axios';
import { API_ROUTES } from '../utils/constants';
import { getTokenFromLocalStorage } from '../lib/auth';

export function Players() {
    const [players, setPlayers] = React.useState([])

    React.useEffect(() => {
        const fetchPlayers = async () => {
            try {
                const response = await axios({
                    method: "GET",
                    url: API_ROUTES.GET_USERS,
                    headers: {
                        authorization: getTokenFromLocalStorage(),
                    }
                })
                setPlayers(response.data)
            }
            catch (err) {
                console.log(err)
            }
            finally {

            }
        }

        fetchPlayers();
    }, [])

    return (
        <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
            {players.map((player, index) => {
                return (<ListItem alignItems="flex-start">
                    <ListItemAvatar>
                        <Avatar alt={player.email} src="/static/images/avatar/1.jpg" />
                    </ListItemAvatar>
                    <ListItemText
                        primary={player.email}
                    />
                </ListItem>
                )
            })}


        </List>
    );
}
// <Divider variant="inset" component="li" />