import { Stack, Tab, TabList, TabPanel, Tabs, Typography, tabClasses } from "@mui/joy";
import React, { useState } from "react";
import { LinkList } from "../LinkList";
import { getOwned, getPopular, getRecent } from "../../services/api/links";
import { NewLinkButton } from "../NewLinkButton";
import { QueryLinkForm } from "../QueryLinkForm";

export const LinkTabs = () => {
    const [lastUpdated, setLastUpdated] = useState(Date.now().toString())

    const onNewLink = () => {
        setLastUpdated(Date.now().toString())
    }

    return (
        <Stack spacing={2}>
            <NewLinkButton onSuccess={onNewLink} />
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
                    <Tab variant="plain" color="neutral"><Typography>Search</Typography></Tab>
                </TabList>
                <TabPanel value={0}><LinkList getLinkFn={getPopular} lastUpdated={lastUpdated} /></TabPanel>
                <TabPanel value={1}><LinkList getLinkFn={getRecent} lastUpdated={lastUpdated} /></TabPanel>
                <TabPanel value={2}><LinkList getLinkFn={getOwned} lastUpdated={lastUpdated} /></TabPanel>
                <TabPanel value={3}><QueryLinkForm /></TabPanel>
            </Tabs>
        </Stack>
    )
}