import { useEffect, useState } from 'react';

import {
    fetchVocabularies
} from '../services/vocabulary.service';

import type { Vocabulary } from '../types/vocabulary';

import VocabularyForm from '../components/vocabulary.form';

const VocabularyPage = () => {

    const [vocabularies, setVocabularies] = useState<Vocabulary[]>([]);

    useEffect(() => {
        const loadVocabularies = async () => {
            try {

                const data =
                    await fetchVocabularies();

                setVocabularies(data);
            } catch(error) {
                console.error(error);
            }
        };

        loadVocabularies();
    },[]);

    return (
        <div>

            <h1>Vocabulary List</h1>

            {vocabularies.map((Vocabulary) => (
                <div key={Vocabulary.id}>
                    <h3>{Vocabulary.word}</h3>
                    <p>{Vocabulary.meaning}</p>
                </div>
            ))}

            <VocabularyForm />
        </div>
    );
};

export default VocabularyPage;