const API_URL = 'http://localhost:8080'
export const API_ROUTES = {
  SIGN_UP: `${API_URL}/users`,
  GET_USERS: `${API_URL}/users`,
  SIGN_IN: `${API_URL}/auth`,
  GET_USER: `${API_URL}/users/me`,
}

export const APP_ROUTES = {
  SIGN_UP: '/sign-up',
  SIGN_IN: '/sign-in',
  SIGN_OUT: '/sign-out',
  GAMES: '/games',
  ONLINE_GAME: '/game',
  PLAYERS: "/players",
  PROFILE: "/profile",
  DASHBOARD: '/',
}