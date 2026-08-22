import * as React from "react";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Box from "@mui/material/Box";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { STRINGS } from "./keys";
import { panelCard } from "./styles";

export default function Navbar() {
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
        mt: 2,
        mb: 1,
        borderRadius: 2,
      }}
    >
      <Tabs
        value={value}
        onChange={handleChange}
        aria-label={STRINGS.nav.ariaLabel}
        centered
        // variant="scrollable"
        scrollButtons="auto"
        sx={{
          minHeight: 52,
          "& .MuiTabs-indicator": {
            backgroundColor: "accent.main",
            height: 3,
            borderRadius: 2,
          },
          "& .MuiTab-root": {
            color: "text.secondary",
            minHeight: 52,
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
    </Box>
  );
}
