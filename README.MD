# ImageBlur
    ImageBlur-go
        Implements pseudo gaussian-blur to an image in Go.
        Uses goroutines to help speed up larger images.
        Displays the image in its own window.
            Updating the image each time a pass is completed.
        By default it averages a 7x7 square of pixels.
            Example:
                With pixel (x, y)
                Runs on x pixels from x-delta to x+delta, and y pixels y-delta to y+delta
                +-------+-------+-------+-------+-------+-------+-------+
                |y-3 x-3|y-3 x-2|y-3 x-1|y-3  x |y-3 x+1|y-3 x+2|y-3 x+3|
                +-------+-------+-------+-------+-------+-------+-------+
                |y-2 x-3|y-2 x-2|y-2 x-1|y-2  x |y-2 x+1|y-2 x+2|y-2 x+3|
                +-------+-------+-------+-------+-------+-------+-------+
                |y-1 x-3|y-1 x-2|y-1 x-1|y-1  x |y-1 x+1|y-1 x+2|y-1 x+3|
                +-------+-------+-------+-------+-------+-------+-------+
                | y  x-3| y  x-2| y  x-1| y   x | y  x+1| y  x+2| y  x+3|
                +-------+-------+-------+-------+-------+-------+-------+
                |y+1 x-3|y+1 x-2|y+1 x-1|y+1  x |y+1 x+1|y+1 x+2|y+1 x+3|
                +-------+-------+-------+-------+-------+-------+-------+
                |y+2 x-3|y+2 x-2|y+2 x-1|y+2  x |y+2 x+1|y+2 x+2|y+2 x+3|
                +-------+-------+-------+-------+-------+-------+-------+
                |y+3 x-3|y+3 x-2|y+3 x-1|y+3  x |y+3 x+1|y+3 x+2|y+3 x+3|
                +-------+-------+-------+-------+-------+-------+-------+

    Usage: /pic <image> <passes>
        To exit press 'Escape' with the image window focused.
        
        If number of passes is not specified, it will run 10 passes.
        Stopping only when 'Escape' is pressed with image window focused.
