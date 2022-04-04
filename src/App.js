import React from "react";

import { createTheme, ThemeProvider, CssBaseline } from "@mui/material";
import Typography from "@mui/material/Typography";

const theme = createTheme({
  typography: {
    fontFamily: "Inter",
  },
});

const App = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Typography variant="h5">Hello world</Typography>
    </ThemeProvider>
  );
};

export default App;
