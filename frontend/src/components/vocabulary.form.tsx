import React, { useState } from "react";

import { createVocabulary } from "../services/vocabulary.service";

const VocabularyForm = () => {
    const [word, setWord] = useState('');

    const [meaning, setMeaning] = useState('');

    const submitHndler = async (
        e: React.SubmitEvent
    ) => {
        e.preventDefault();

        try {
            await createVocabulary({
                word,
                meaning,
                partOfSpeech: 'noun'
            });

            setWord('');
            setMeaning('');

            alert('Created');
        } catch(error) {
            console.error(error);
        }
    };

    return (
        <form onSubmit={submitHndler}>

            <input
                type="text"
                placeholder="word"
                value={word}
                onChange={(e) =>
                    setWord(e.target.value)
                }
             />

             <input
                type="text"
                placeholder="meaning"
                value={meaning}
                onChange={(e) =>
                    setMeaning(e.target.value)
                }
             />

             <button type="submit">
                Create
             </button>
        </form>
    );
};

export default VocabularyForm;