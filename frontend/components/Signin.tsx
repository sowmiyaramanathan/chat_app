import { Alert, Button, IconButton, InputAdornment, Stack } from "@mui/material";
import axios from "axios";
import { Formik } from "formik";
import { useRouter } from "next/router";
import CancelIcon from "@mui/icons-material/Cancel";
import LoginIcon from "@mui/icons-material/Login";
import { useState } from "react";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import * as Yup from "yup";
import { setToken, setPvtKey } from "../token/token";
import { CustomTextField } from "./CustomComponets";
import { STRINGS } from "./keys";
import { containedButton, outlinedButton, panelCard } from "./styles";
import { getApiErrorCode, getApiErrorMessage } from "./api";

export default function Signin() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const handleClickShowPassword = () => setShowPassword((prev) => !prev);

  const validationSchema = Yup.object().shape({
    username: Yup.string().required(STRINGS.validation.usernameRequired),
    password: Yup.string()
      .min(8, STRINGS.validation.passwordMin)
      .required(STRINGS.validation.passwordRequired),
  });

  return (
    <Formik
      initialValues={{ username: "", password: "" }}
      validationSchema={validationSchema}
      onSubmit={async (values, { setFieldError }) => {
        setSubmitError(null);
        try {
          const response = await axios.post("http://localhost:8000/user/login", {
            Username: values.username,
            Password: values.password,
          });
          setToken(response.data.token);
          setPvtKey(response.data.privateKey);
          await router.push("/user/profile");
        } catch (error: unknown) {
          const code = getApiErrorCode(error);
          if (code === "username") {
            setFieldError("username", STRINGS.errors.usernameNotFound);
          } else if (code === "password") {
            setFieldError("password", STRINGS.errors.wrongPassword);
          } else {
            setSubmitError(getApiErrorMessage(error, STRINGS.errors.signIn));
          }
        }
      }}
    >
      {({ values, errors, touched, handleChange, handleSubmit }) => {
        return (
          <form onSubmit={handleSubmit}>
            <Stack
              sx={{
                ...panelCard,
                gap: 2,
                maxWidth: 420,
                width: "100%",
                mx: "auto",
                p: { xs: 3, sm: 4 },
                mt: { xs: 4, md: 8 },
              }}
            >
              {submitError && <Alert severity="error">{submitError}</Alert>}
              <CustomTextField
                id="username"
                label={STRINGS.auth.username}
                value={values.username}
                variant="outlined"
                type="text"
                onChange={handleChange}
                error={touched.username && Boolean(errors.username)}
                helperText={touched.username && errors.username}
              />

              <CustomTextField
                id="password"
                label={STRINGS.auth.password}
                value={values.password}
                variant="outlined"
                type={showPassword ? "text" : "password"}
                onChange={handleChange}
                error={touched.password && Boolean(errors.password)}
                helperText={touched.password && errors.password}
                slotProps={{
                  input: {
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          aria-label={STRINGS.auth.togglePasswordVisibility}
                          onClick={handleClickShowPassword}
                          edge="end"
                          sx={{ color: "primary.main" }}
                        >
                          {showPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    ),
                  },
                }}
              />
              <Stack direction="row" sx={{ justifyContent: "space-between", gap: 2, pt: 1 }}>
                <Button
                  variant="outlined"
                  sx={outlinedButton}
                  startIcon={<CancelIcon />}
                  onClick={() => router.back()}
                >
                  {STRINGS.auth.cancel}
                </Button>
                <Button
                  variant="contained"
                  type="submit"
                  sx={containedButton}
                  startIcon={<LoginIcon />}
                >
                  {STRINGS.auth.signIn}
                </Button>
              </Stack>
            </Stack>
          </form>
        );
      }}
    </Formik>
  );
}
