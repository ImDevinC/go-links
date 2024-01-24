import { Box, Stack } from "@mui/joy";
import React from "react";
import { QueryLinkForm } from "../QueryLinkForm";
import { LinkTabs } from "../LinkTabs";

export const LinkOverview = () => {
    return (
        <Box sx={{ flex: 1, width: '100%', mx: 'auto' }}>
            <Stack spacing={4} sx={{
                display: 'flex',
                maxWidth: '1600px',
                mx: 'auto',
                px: { xs: 2, md: 6 },
                py: { xs: 2, md: 3 }
            }}>
                <QueryLinkForm />
                <LinkTabs />
            </Stack>
        </Box>
    )
}