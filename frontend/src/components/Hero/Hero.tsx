import { Box, Stack, Typography } from "@mui/joy";
import React from "react";

export const Hero = () => {
    return (
        <Box sx={{ flex: 1, width: '100%', mx: 'auto' }}>
            <Stack spacing={4} sx={{
                display: 'flex',
                maxWidth: '1600px',
                mx: 'auto',
            }}>
                <Typography level="h1">go/</Typography>
                <Typography level="body-md">Internal URL shortener to provide quick and memorable URLs.</Typography>
            </Stack>
        </Box>

    )
}