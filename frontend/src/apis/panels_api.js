import axios from 'axios';

export const getAllPanels = () => axios.get("/api/panels")
