import { Box, Stack } from "@mui/joy"
import { LinkOverview } from "../../components/LinkOverview"
import { ThemeToggle } from "../../components/ThemeToggle"
import { Hero } from "../../components/Hero"

export const Home = () => {
    return (
        <Box sx={{ flex: 1, width: '100%' }}>
            <Stack spacing={2} pt={2} pr={2}>
                <ThemeToggle />
                <Hero />
                <LinkOverview />
            </Stack>
        </Box>
    )
}