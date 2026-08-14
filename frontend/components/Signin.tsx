import { Button, IconButton, InputAdornment, Stack } from "@mui/material";
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

export default function Signin() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);

  const handleClickShowPassword = () => setShowPassword((prev) => !prev);

  const validationSchema = Yup.object().shape({
    username: Yup.string().required("Username is required"),
    password: Yup.string()
      .min(8, "Enter minimum 8 characters")
      .required("Password is required"),
  });

  return (
    <Formik
      initialValues={{ username: "", password: "" }}
      validationSchema={validationSchema}
      onSubmit={async (values, { setFieldError }) => {
        axios({
          method: "post",
          url: "http://localhost:8000/user/login",
          data: {
            Username: values.username,
            Password: values.password,
          },
        })
          .then((response) => {
            setToken(response.data.token);
            setPvtKey(response.data.privateKey);
            router.push("/user/profile");
          })
          .catch((error) => {
            const message = error.response.data.message;
            if (message == "Username") {
              setFieldError("username", "Username does not exist");
            } else if (message == "Password") {
              setFieldError("password", "Wrong Password. Try again");
            } else {
              console.log(error);
            }
          });
      }}
    >
      {({ values, errors, touched, handleChange, handleSubmit }) => {
        return (
          <form onSubmit={handleSubmit}>
            <Stack sx={{ gap: 2, maxWidth: "50vh", margin: "auto", pt: "20vh" }}>
              <CustomTextField
                id="username"
                label="Username"
                value={values.username}
                variant="outlined"
                type="text"
                onChange={handleChange}
                error={touched.username && Boolean(errors.username)}
                helperText={touched.username && errors.username}
              />

              <CustomTextField
                id="password"
                label="Password"
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
                          aria-label="toggle password visibility"
                          onClick={handleClickShowPassword}
                          edge="end"
                          sx={{ color: "msgBg.main" }}
                        >
                          {showPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    ),
                  },
                }}
              />
              <Stack direction="row" sx={{justifyContent: "space-around", gap: 2}}>
                <Button
                  variant="outlined"
                  color="primary"
                  sx={{
                    color: "primary.light",
                    borderColor: "msgBg.main",
                    ":hover": {
                      color: "primary.contrastText",
                      borderColor: "primary.light",
                    },
                  }}
                  startIcon={<CancelIcon />}
                  onClick={() => router.back()}
                >
                  Cancel
                </Button>
                <Button
                  variant="contained"
                  type="submit"
                  sx={{
                    color: "primary.contrastText",
                    backgroundColor: "primary.main",
                    ":hover": {
                      color: "primary.light",
                    },
                  }}
                  startIcon={<LoginIcon />}
                >
                  Sign In
                </Button>
              </Stack>
            </Stack>
          </form>
        );
      }}
    </Formik>
  );
}
