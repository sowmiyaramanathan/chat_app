import { Box, Button, Typography } from "@mui/material";
import { useRouter } from "next/router";

export default function Home() {
  const router = useRouter();

  function signup() {
    router.push({
      pathname: "/user/signup",
    });
  }

  return (
    <Box
      sx={{
        display: "flex",
        height: "80vh",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
      }}
    >
      <Typography color="primary.light" sx={{ mb: 2 }}>
        Mini Chat Application
      </Typography>
      <Button
        variant="contained"
        onClick={signup}
        sx={{
          color: "primary.contrastText",
          ":hover": {
            backgroundColor: "primary.main",
            color: "primary.light",
          },
        }}
      >
        Signup Now
      </Button>
    </Box>
  );
}
