import * as React from "react";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import LightModeRoundedIcon from "@mui/icons-material/LightModeRounded";
import DarkModeRoundedIcon from "@mui/icons-material/DarkModeRounded";
import AutoAwesomeRoundedIcon from "@mui/icons-material/AutoAwesomeRounded";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { STRINGS } from "./keys";
import { panelCard } from "./styles";

export default function Navbar({
  mode,
  onToggleMode,
}: {
  mode: "light" | "dark";
  onToggleMode: () => void;
}) {
  const [hasToken, setHasToken] = useState(false);
  const [value, setValue] = useState(0);
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem("token");
    setHasToken(Boolean(token));

    const paths = token
      ? ["/user/profile", "/user/chats", "/user/requests"]
      : ["/", "/user/signup", "/user/signin"];

    const currentIndex = router.pathname.startsWith("/user/chats")
      ? 1
      : paths.indexOf(router.pathname);
    setValue(currentIndex >= 0 ? currentIndex : 0);
  }, [router.pathname]);

  const handleChange = (_event: React.SyntheticEvent, newValue: number) => {
    const paths = hasToken
      ? ["/user/profile", "/user/chats", "/user/requests"]
      : ["/", "/user/signup", "/user/signin"];

    setValue(newValue);
    router.push(paths[newValue]);
  };

  return (
    <Box
      sx={{
        ...panelCard,
        mx: { xs: 2, sm: 3 },
        mt: { xs: 1.5, sm: 2.5 },
        mb: 1,
        borderRadius: 4,
        px: { xs: 1, sm: 2 },
        display: "flex",
        alignItems: "center",
        gap: 1,
      }}
    >
      <Box sx={{ display: "flex", alignItems: "center", gap: 0.75, pl: 1, minWidth: "fit-content" }}>
        <AutoAwesomeRoundedIcon sx={{ color: "accent.main", fontSize: 22 }} />
        <Typography sx={{ display: { xs: "none", sm: "block" }, fontWeight: 700, color: "text.primary", letterSpacing: "0.02em" }}>
          {STRINGS.app.name}
        </Typography>
      </Box>
      <Tabs
        value={value}
        onChange={handleChange}
        aria-label={STRINGS.nav.ariaLabel}
        variant="scrollable"
        scrollButtons="auto"
        slotProps={{ indicator: { sx: { backgroundColor: "accent.main", height: 4, borderRadius: 4 } } }}
        sx={{
          minHeight: 52,
          flex: 1,
          "& .MuiTab-root": {
            color: "text.secondary",
            minHeight: 52,
            minWidth: { xs: 76, sm: 96 },
            fontSize: { xs: "0.84rem", sm: "0.95rem" },
            "&.Mui-selected": {
              color: "primary.main",
            },
          },
        }}
      >
        {hasToken && <Tab label={STRINGS.nav.profile} />}
        {hasToken && <Tab label={STRINGS.nav.chats} />}
        {hasToken && <Tab label={STRINGS.nav.requests} />}

        {!hasToken && <Tab label={STRINGS.nav.home} />}
        {!hasToken && <Tab label={STRINGS.nav.signUp} />}
        {!hasToken && <Tab label={STRINGS.nav.signIn} />}
      </Tabs>
      <Tooltip title={mode === "light" ? "Use night mode" : "Use day mode"}>
        <IconButton
          onClick={onToggleMode}
          aria-label={mode === "light" ? "Use night mode" : "Use day mode"}
          sx={{ color: "primary.main", bgcolor: "surface.dark", mr: 0.25 }}
        >
          {mode === "light" ? <DarkModeRoundedIcon /> : <LightModeRoundedIcon />}
        </IconButton>
      </Tooltip>
    </Box>
  );
}
