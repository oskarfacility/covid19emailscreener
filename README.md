# E-Mail Categorizer

Categorize E-Mails based on predefined rules.

As many hospitals and health organisations are facing a shortage of personell 
because of the COVID-19/Corona crisis and publicly asked former nurses, 
doctors, and rescue worker for help, some organisations received thousands 
of e-mails.

This tool can help to categorize the e-mails based on their content.
The tool is not perfect, but is doing a good job already at some organizations.

## How to use it?

1. Export the E-Mails into a folder called `emails`, they should be in the 
`.eml` format
2. Drop the `EmailCategorizer.exe` (WIN) or `EmailCategorizer_macos` (Mac) at 
the same level as the `emails` folder
3. Drop a `config.yml` at the same level as `emails` and the executable

It should look like this

```
 > emails/
 config.yml
 EmailCategorizer*
```

Execute the programm, a `.csv` should be generated.
