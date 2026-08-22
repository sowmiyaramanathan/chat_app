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
import { STRINGS } from "./keys";
import { containedButton, outlinedButton, panelCard } from "./styles";

export default function Signup() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const handleClickShowPassword = () => setShowPassword((prev) => !prev);
  const handleClickShowConfirmPassword = () =>
    setShowConfirmPassword((prev) => !prev);

  const validationSchema = Yup.object().shape({
    name: Yup.string().required(STRINGS.validation.nameRequired),
    username: Yup.string().required(STRINGS.validation.usernameRequired),
    mobile_number: Yup.string()
      .matches(/^[6-9]\d{9}$/, {
        message: STRINGS.validation.mobileInvalid,
      })
      .required(STRINGS.validation.mobileRequired),
    password: Yup.string()
      .min(8, STRINGS.validation.passwordMin)
      .required(STRINGS.validation.passwordRequired),
    confirm_password: Yup.string()
      .oneOf([Yup.ref("password")], STRINGS.validation.passwordsMustMatch)
      .required(STRINGS.validation.confirmPasswordRequired),
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
              setFieldError("username", STRINGS.errors.usernameExists);
            } else if (message == "Number") {
              setFieldError("mobile_number", STRINGS.errors.mobileExists);
            } else {
              console.log(error);
            }
          });
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
                mt: { xs: 4, md: 6 },
              }}
            >
              <CustomTextField
                id="name"
                label={STRINGS.auth.name}
                value={values.name}
                variant="outlined"
                onChange={handleChange}
                error={touched.name && Boolean(errors.name)}
                helperText={touched.name && errors.name}
              />
              <CustomTextField
                id="username"
                label={STRINGS.auth.username}
                value={values.username}
                variant="outlined"
                onChange={handleChange}
                error={touched.username && Boolean(errors.username)}
                helperText={touched.username && errors.username}
              />
              <CustomTextField
                id="mobile_number"
                label={STRINGS.auth.mobileNumber}
                value={values.mobile_number}
                variant="outlined"
                onChange={handleChange}
                error={touched.mobile_number && Boolean(errors.mobile_number)}
                helperText={touched.mobile_number && errors.mobile_number}
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
              <CustomTextField
                id="confirm_password"
                label={STRINGS.auth.confirmPassword}
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
                          aria-label={STRINGS.auth.togglePasswordVisibility}
                          onClick={handleClickShowConfirmPassword}
                          edge="end"
                          sx={{ color: "primary.main" }}
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
                  startIcon={<HowToRegIcon />}
                >
                  {STRINGS.auth.signUp}
                </Button>
              </Stack>
            </Stack>
          </form>
        );
      }}
    </Formik>
  );
}
