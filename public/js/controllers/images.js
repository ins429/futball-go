app.controller('Images', function ImagesCtrl($scope, $route, $routeParams, api) {
  $scope.initImages = function() {
    $scope.images = [];
    $scope.getMoreImages();
    $scope.imageData = {};
    $scope.showImageForm = false;
    $scope.formUrl = '/views/images/form.html';
  };

  $scope.getMoreImages = function() {
    api.queryImages().then(function(data) {
      if (data && data.response) {
        $scope.images = data.response
      } else {
        // fix me, handle error
      }
    });
  };

  // pop form to create an image
  $scope.newImage = function() {
    $scope.imageData = {
      action: 'New'
    };
    $scope.togglerImageForm();
  };

  // pop form to edit an image
  $scope.editImage = function(image) {
    $scope.imageData = {
      action: 'Edit',
      modalName: image.name,
      id: image.id,
      name: image.name,
      tags: image.tags,
      url: image.url
    }
    $scope.togglerImageForm();
  };

  $scope.togglerImageForm = function() {
    $scope.showImageForm = !$scope.showImageForm;
    if (angular.element('#submit-button').hasClass('in-progress')) {
      angular.element('#submit-button').html('<i class="fa fa-cloud-upload"></i>Submit').removeClass('in-progress');
    }
  };

  $scope.submitImage = function(imageData) {
    angular.element('#submit-button').html('<i class="fa fa-cloud-upload"></i>Uploading').addClass('in-progress');
    if (imageData.action.toLowerCase() === 'edit') {
      updateImage(imageData)
    } else {
      createImage(imageData)
    }
  };

  $scope.deleteImage = function(imageId) {
    if (confirm("Do you want to delete this image?")) {
      api.deleteImage({
        id: imageId
      }).then(function(data) {
        if (data) {
          deleteImageFromTheList(imageId)
        } else {
          // fix me, handle error
        }
      });
    }
  };

  $scope.searchImages = function(query) {
    api.queryImages({
      query: query
    }).then(function(data) {
      if (data && data.response) {
        console.log(data)
        $scope.images = data.response
      } else {
        // fix me, handle error
      }
    });
  }

  $scope.onFileSelect = function($files) {
    for (var i = 0; i < $files.length; i++) {
      $scope.imageData.file = $files[i];
    }
  }

  function updateImage(imageData) {
    var params = buildImageData(imageData);
    if ($scope.imageData.file) {
      api.updateImageWithFile(params, imageData.file).then(function(resp) {
        if (resp && resp.status && resp.status === 200) {
          updateImageFromTheList(imageData.id, resp.data.response[0]);
          $scope.togglerImageForm();
        } else {
          // fix me, handle error
        }
        angular.element('#submit-button').html('<i class="fa fa-cloud-upload"></i>Submit').removeClass('in-progress');
        $scope.imageData = {};
      });
    } else {
      api.updateImage(params).then(function(resp) {
        if (resp && resp.status && resp.status === 200) {
          console.log(resp)
          updateImageFromTheList(imageData.id, resp.response[0]);
          $scope.togglerImageForm();
        } else {
          // fix me, handle error
        }
        angular.element('#submit-button').html('<i class="fa fa-cloud-upload"></i>Submit').removeClass('in-progress');
        $scope.imageData = {};
      });
    }
  };

  function createImage(imageData) {
    var params = buildImageData(imageData);
    if ($scope.imageData.file) {
      api.createImageWithFile(params, imageData.file).then(function(resp) {
        if (resp && resp.status && resp.status === 200) {
          $scope.images.unshift(resp.data.response[0])
          $scope.togglerImageForm();
        } else {
          // fix me, handle error
        }
        angular.element('#submit-button').html('<i class="fa fa-cloud-upload"></i>Submit').removeClass('in-progress');
        $scope.imageData = {};
      });
    } else {
      api.createImage(params).then(function(data) {
        if (data && data.status && data.status === 200) {
          $scope.images.unshift(data.response[0])
          $scope.togglerImageForm();
        } else {
          // fix me, handle error
        }
        angular.element('#submit-button').html('<i class="fa fa-cloud-upload"></i>Submit').removeClass('in-progress');
        $scope.imageData = {};
      });
    }
  };

  function buildImageData(imageData) {
    var data = {
      name: imageData.name
    }
    if (imageData.id) {
      data.id = imageData.id
    };
    if (imageData.url && !imageData.file) {
      data.url = imageData.url
    };
    if (imageData.tags) {
      data.tags = imageData.tags
    };
    return data
  }

  function updateImageFromTheList(imageId, data) {
    if (!$scope.images && $scope.images.length === 0) {
      return;
    } else {
      for (var i=0; i < $scope.images.length; i++) {
        angular.forEach($scope.images[i], function(value, key){
          if (key === 'id' && value === imageId) {
            $scope.images[i] = data;
            return;
          }
        });
      }
    }
  }

  function deleteImageFromTheList(imageId) {
    if (!$scope.images && $scope.images.length === 0) {
      return;
    } else {
      for (var i=0; i < $scope.images.length; i++) {
        angular.forEach($scope.images[i], function(value, key){
          if (key === 'id' && value === imageId) {
            $scope.images.splice(i, 1);
            return;
          }
        });
      }
    }
  }
});