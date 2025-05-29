import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import Todo from "./Todo";
import {
  Box, Typography, Menu, MenuItem, IconButton,
} from '@mui/material';
import SettingsIcon from '@mui/icons-material/Settings';

function Dashboard() {
  const navigate = useNavigate();
  const token = localStorage.getItem("token");
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  useEffect(() => {
    if (!token) navigate("/");
  }, [navigate, token]);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  const handleProfileUpdate = () => {
    navigate("/profile");
  };

  const handleMenuClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <Box className="container">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" sx={{ fontWeight: 'bold', color: '#ffffff' }}>
          Todo List
        </Typography>
        <IconButton onClick={handleMenuClick} sx={{ color: '#ffffff' }}>
          <SettingsIcon />
        </IconButton>
      </Box>

      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <MenuItem onClick={() => { handleClose(); handleProfileUpdate(); }}>
          Profilni yangilash
        </MenuItem>
        <MenuItem onClick={() => { handleClose(); handleLogout(); }}>
          Chiqish
        </MenuItem>
      </Menu>

      <Todo />
    </Box>
  );
}

export default Dashboard;
