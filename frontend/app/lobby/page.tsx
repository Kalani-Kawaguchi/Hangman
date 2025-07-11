'use client';
import { Suspense } from 'react';
import Lobby from './lobby';

export default function Page() {
    return(
        <Suspense fallback={<div>Loading lobby...</div>}>
            <Lobby />
        </Suspense>
    )
}