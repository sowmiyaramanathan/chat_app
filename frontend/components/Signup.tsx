import { Button, IconButton, InputAdornment, Stack } from "@mui/material";
import axios from "axios";
import { Formik } from "formik";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import { useRouter } from "next/router";
import { useState } from "react";
import HowToRegIcon from "@mui/icons-material/HowToReg";
import CancelIcon from "@mui/icons-material/Cancel";
import * as Yup from "yup";
import { CustomTextField } from "./CustomComponets";

export default function Signup() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const handleClickShowPassword = () => setShowPassword((prev) => !prev);
  const handleClickShowConfirmPassword = () =>
    setShowConfirmPassword((prev) => !prev);

  const validationSchema = Yup.object().shape({
    name: Yup.string().required("Name is required"),
    username: Yup.string().required("Username is required"),
    mobile_number: Yup.string()
      .matches(/^[6-9]\d{9}$/, {
        message: "Please enter valid number.",
      })
      .required("Mobile number is required"),
    password: Yup.string()
      .min(8, "Enter minimum 8 characters")
      .required("Password is required"),
    confirm_password: Yup.string()
      .oneOf([Yup.ref("password")], "Passwords must match")
      .required("Confirm Password is required"),
  });

  return (
    <Formik
      initialValues={{
        name: "",
        username: "",
        mobile_number: "",
        password: "",
        confirm_password: "",
      }}
      validationSchema={validationSchema}
      onSubmit={async (values, { setFieldError }) => {
        await axios({
          method: "post",
          url: "http://localhost:8000/user/register",
          data: {
            Name: values.name,
            Username: values.username,
            mobile_number: values.mobile_number,
            Password: values.password,
          },
        })
          .then(() => {
            router.push("/user/signin");
          })
          .catch((error) => {
            const message = error.response.data.message;
            if (message == "Username") {
              console.log("In");
              setFieldError("username", "Username already exists");
            } else if (message == "Number") {
              setFieldError("mobile_number", "Mobile number already exists");
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
                id="name"
                label="Name"
                value={values.name}
                variant="outlined"
                onChange={handleChange}
                error={touched.name && Boolean(errors.name)}
                helperText={touched.name && errors.name}
              />
              <CustomTextField
                id="username"
                label="Username"
                value={values.username}
                variant="outlined"
                onChange={handleChange}
                error={touched.username && Boolean(errors.username)}
                helperText={touched.username && errors.username}
              />
              <CustomTextField
                id="mobile_number"
                label="Mobile Number"
                value={values.mobile_number}
                variant="outlined"
                onChange={handleChange}
                error={touched.mobile_number && Boolean(errors.mobile_number)}
                helperText={touched.mobile_number && errors.mobile_number}
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
              <CustomTextField
                id="confirm_password"
                label="Confirm Password"
                value={values.confirm_password}
                type={showConfirmPassword ? "text" : "password"}
                variant="outlined"
                onChange={handleChange}
                error={
                  touched.confirm_password && Boolean(errors.confirm_password)
                }
                helperText={touched.confirm_password && errors.confirm_password}
                slotProps={{
                  input: {
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          aria-label="toggle password visibility"
                          onClick={handleClickShowConfirmPassword}
                          edge="end"
                          sx={{ color: "msgBg.main" }}
                        >
                          {showConfirmPassword ? (
                            <VisibilityOff />
                          ) : (
                            <Visibility />
                          )}
                        </IconButton>
                      </InputAdornment>
                    ),
                  },
                }}
              />
              <Stack direction="row" sx={{justifyContent: "space-evenly"}}  spacing={2}>
                <Button
                  variant="outlined"
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
                  startIcon={<HowToRegIcon />}
                >
                  Sign up
                </Button>
              </Stack>
            </Stack>
          </form>
        );
      }}
    </Formik>
  );
}
