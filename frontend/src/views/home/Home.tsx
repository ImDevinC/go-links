import { Box, Stack } from "@mui/joy"
import { LinkOverview } from "../../components/LinkOverview"
import { ThemeToggle } from "../../components/ThemeToggle"

export const Home = () => {
    return (
        <Box sx={{ flex: 1, width: '100%' }}>
            <Stack spacing={2} pt={2} pr={2}>
                <ThemeToggle />
                <LinkOverview />
            </Stack>
        </Box>
    )
}