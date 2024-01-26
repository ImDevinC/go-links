import { Box } from "@mui/joy"
import { LinkOverview } from "../../components/LinkOverview"
import { ThemeToggle } from "../../components/ThemeToggle"

export const Home = () => {
    return (
        <Box sx={{ flex: 1, width: '100%' }}>
            <ThemeToggle />
            <LinkOverview />
        </Box>
    )
}