import { Box, Tab, TabList, TabPanel, Tabs, Typography, tabClasses } from "@mui/joy";
import React from "react";
import { PopularLinks } from "../PopularLinks";

export const LinkTabs = () => {
    return (
        <Box sx={{ flex: 1, maxWidth: '800', mx: 'auto' }}>
            <Tabs>
                <TabList sx={{
                    pt: 1,
                    justifyContent: 'center',
                    [`&& .${tabClasses.root}`]: {
                        flex: 'initial',
                        bgcolor: 'transparent',
                        '&:hover': {
                            bgcolor: 'transparent',
                        }
                    },
                    [`&.${tabClasses.selected}`]: {
                        color: 'primary.plainColor',
                        '&::after': {
                            height: 2,
                            borderTopLeftRadius: 3,
                            borderTopRightRadius: 3,
                            bgcolor: 'primary.500',
                        },
                    }
                }} >
                    <Tab variant="plain" color="neutral"><Typography>Popular</Typography></Tab>
                    <Tab variant="plain" color="neutral"><Typography>Recent</Typography></Tab>
                    <Tab variant="plain" color="neutral"><Typography>My Links</Typography></Tab>
                    <Tab variant="plain" color="neutral"><Typography>Recently deactivated</Typography></Tab>
                </TabList>
                <TabPanel value={0}><PopularLinks /></TabPanel>
            </Tabs>
        </Box >
    )
}