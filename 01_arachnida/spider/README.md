# About Spider
The spider prograp scrapes a website recursively looking for images.
The first URL passed will be the entrypoint. With no flag spider will scrape for images an download them.
With the recursive flag each URL found, belonging to the same host as the Arg (entrypoint), will be scraped. 

## Usage

`./spider [-rlp] URL`
- Option -r : recursively downloads the images in a URL received as a parameter.
- Option -r -l [N] : indicates the maximum depth level of the recursive download.
If not indicated, it will be 5.
- Option -p [PATH] : indicates the path where the downloaded files will be saved.
If not specified, `./data/` will be used.
- Option -v : verbosely display each step
- Option --spread : scrape each found URL, regardless of Host

The program will download the following extensions by default:
- .jpg/jpeg
- .png
- .gif
- .bmp
