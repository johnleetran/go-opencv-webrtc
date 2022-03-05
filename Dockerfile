FROM ubuntu:18.04 

#install some image libs
WORKDIR /tmp
RUN apt-get update 
RUN apt-get -y install curl libpng-dev libjpeg-dev libtiff-dev libwebp-dev zlib1g-dev libvpx-dev pkg-config git curl zip unzip tar ninja-build

#install clang
WORKDIR /tmp
RUN apt-get -y install wget lsb-core software-properties-common build-essential
RUN wget https://apt.llvm.org/llvm.sh
RUN chmod +x llvm.sh
RUN ./llvm.sh 11

#install cmake
WORKDIR /tmp
RUN wget -qO- "https://cmake.org/files/v3.20/cmake-3.20.0-linux-x86_64.tar.gz" | tar --strip-components=1 -xz -C /usr/local

#install opencv core and contrib
ENV OPENCV_VERSION=4.5.5
WORKDIR /opt
RUN curl -L https://github.com/opencv/opencv/archive/refs/tags/${OPENCV_VERSION}.zip -o opencv_${OPENCV_VERSION}.zip
RUN curl -L https://github.com/opencv/opencv_contrib/archive/refs/tags/${OPENCV_VERSION}.zip -o opencv_contrib_${OPENCV_VERSION}.zip
RUN unzip opencv_${OPENCV_VERSION}.zip
RUN unzip opencv_contrib_${OPENCV_VERSION}.zip
WORKDIR /opt/opencv-4.5.5/build
RUN cmake -GNinja -D OPENCV_GENERATE_PKGCONFIG=YES -DOPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-4.5.5/modules ..
RUN ninja
RUN ninja install 

# install golang
WORKDIR /opt 
RUN curl -L https://go.dev/dl/go1.17.8.linux-amd64.tar.gz -O 
RUN gunzip -c go1.17.8.linux-amd64.tar.gz | tar xv 
RUN mv go /usr/local/go
ENV PATH=${PATH}:/usr/local/go/bin

WORKDIR /app
ENV CGO_CXXFLAGS="--std=c++11"
ENV LD_LIBRARY_PATH="/opt/opencv-4.5.5/build/lib"
ENV CGO_LDFLAGS="-L/opt/opencv-4.5.5/build/lib -lopencv_stitching -lopencv_superres -lopencv_videostab -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_dnn_objdetect -lopencv_dpm -lopencv_face -lopencv_photo -lopencv_fuzzy -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_optflow -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_stereo -lopencv_structured_light -lopencv_phase_unwrapping -lopencv_surface_matching -lopencv_tracking -lopencv_datasets -lopencv_dnn -lopencv_plot -lopencv_xfeatures2d -lopencv_shape -lopencv_video -lopencv_ml -lopencv_ximgproc -lopencv_calib3d -lopencv_features2d -lopencv_highgui -lopencv_videoio -lopencv_flann -lopencv_xobjdetect -lopencv_imgcodecs -lopencv_objdetect -lopencv_xphoto -lopencv_imgproc -lopencv_core"
ADD . .
RUN go build main.go
CMD ["bash", "-l", "-c", "./main"]