import { TextField, styled } from "@mui/material";

const CustomTextField = styled(TextField)(({ theme }) => ({
  "& .MuiFormLabel-root": {
    color: theme.palette.text.secondary,
    "&.Mui-focused": {
      color: theme.palette.primary.main,
    },
  },
  "& .MuiOutlinedInput-notchedOutline": {
    borderColor: theme.palette.secondary.dark,
  },
  "& .MuiOutlinedInput-root": {
    bgcolor: theme.palette.surface.light,
    borderRadius: 12,
    "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
      borderColor: theme.palette.primary.main,
    },
    "&:hover .MuiOutlinedInput-notchedOutline": {
      borderColor: theme.palette.primary.main,
    },
  },
  "& .MuiInputBase-input": {
    color: theme.palette.text.primary,
  },
  "& .MuiFormHelperText-root": {
    color: theme.palette.error.main,
  },
}));

const MessageField = styled(TextField)(({ theme }) => ({
  "& .MuiInputBase-input::placeholder": {
    color: theme.palette.text.secondary,
    opacity: 0.8,
  },
  "& .MuiOutlinedInput-notchedOutline": {
    borderColor: theme.palette.secondary.dark,
  },
  "& .MuiOutlinedInput-root": {
    bgcolor: theme.palette.surface.light,
    borderRadius: 12,
    "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
      borderColor: theme.palette.primary.main,
    },
    "&:hover .MuiOutlinedInput-notchedOutline": {
      borderColor: theme.palette.primary.main,
    },
  },
  "& .MuiInputBase-input": {
    color: theme.palette.text.primary,
  },
}));

export { CustomTextField, MessageField };
