# GopherEx

## Spec requirements
- Text
    - Custom formating
- Forms
    - Custom header containing form information
    - Multiple forms of input
- Images
- Videos
- File download

## GopherEx struct
    header {
        key:val
    }

    style_init_byte [
        (style_byte args*)
    ]

    generator_byte args* [
        (type_byte content*)*
    ]

    type_byte [
        (type_byte content*)*
    ]



    "generator_byte args [
        type_byte content
    ]"