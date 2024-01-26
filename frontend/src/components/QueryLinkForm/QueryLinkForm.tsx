import { FormControl } from "@mui/base";
import { Alert, Button, Input, Stack } from "@mui/joy";
import React, { useState } from "react";
import { LinkList } from "../LinkList";
import { LinkData, queryLinks } from "../../services/api/links";

export const QueryLinkForm = () => {
    const [query, setQuery] = useState('');
    const [links, setLinks] = useState<LinkData[]>([])
    const [errorMessage, setErrorMessage] = useState('')

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault()
        setErrorMessage('')
        search(query)
    }

    const search = async (query: string) => {
        const result = await queryLinks(query)
        if (result.error) {
            setErrorMessage(result.error)
            return
        }
        setLinks(result.links)
    }

    return (
        <Stack spacing={2}>
            <form onSubmit={handleSubmit} id="find-link">
                <FormControl
                    required
                    value={query}
                >
                    <Input
                        sx={{ '--Input-decoratorChildHeight': '45px' }}
                        placeholder="Search name, description or URL"
                        name="query"
                        onChange={(event) => { setQuery(event.target.value) }}
                        endDecorator={
                            <Button type="submit" color="primary" sx={{ borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}>Search</Button>
                        }
                    />
                </FormControl>
            </form>
            {errorMessage !== '' && <Alert color="danger">{errorMessage}</Alert>}
            <LinkList getLinkFn={async () => { return { links: links } }} lastUpdated={Date.now().toString()} />
        </Stack>
    )
}