import { FormControl } from "@mui/base";
import { Button, Input, Stack } from "@mui/joy";
import React, { useState } from "react";
import { NewLinkButton } from "../NewLinkButton";

export const QueryLinkForm = () => {
    const [query, setQuery] = useState('');
    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        // TODO: implement
    }

    return (
        <Stack spacing={2}>
            <NewLinkButton />
            <form onSubmit={handleSubmit} id="find-link">
                <FormControl>
                    <Input
                        sx={{ '--Input-decoratorChildHeight': '45px' }}
                        placeholder="Search name, description or URL"
                        name="query"
                        required
                        value={query}
                        onChange={(event) => { setQuery(event.target.value) }}
                        endDecorator={
                            <Button type="submit" color="primary" sx={{ borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}>Search</Button>
                        }
                    />
                </FormControl>
            </form>
        </Stack>
    )
}